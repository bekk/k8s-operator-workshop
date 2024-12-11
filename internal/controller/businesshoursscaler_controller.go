package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bekk/k8s-operator-workshop/internal/clock"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/bekk/k8s-operator-workshop/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
)

// Important bits below - RBAC rules needed for the controller to function properly. These annotations
// are used by the controller-runtime library to generate the necessary RBAC rules for the controller.

// +kubebuilder:rbac:groups=apps.k8s.bekk.no,resources=businesshoursscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.k8s.bekk.no,resources=businesshoursscalers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.k8s.bekk.no,resources=businesshoursscalers/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update;patch

const (
	fakeTimeAnnotationKey     = "bhs.bekk.no/fake-time"
	intervalBetweenReconciles = 10 * time.Second

	ScaleSuccess = "DeploymentScaled"
	ScaleFailure = "DeploymentScaleFailed"
)

// Construct our reconciler struct with required fields
type BusinessHoursScalerReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Clock    clock.Clock
}

// SetupWithManager sets up the controller with the Manager.
func (r *BusinessHoursScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.BusinessHoursScaler{}).
		Complete(r)
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the BusinessHoursScaler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *BusinessHoursScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling instance")

	// Fetch the BusinessHoursScaler resource
	bhs := &appsv1alpha1.BusinessHoursScaler{}
	if err := r.Get(ctx, req.NamespacedName, bhs); err != nil {
		if kerrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Parse start and end times with time zone
	startTime, validStart := parseTime(bhs.Spec.StartTime, bhs.Spec.TimeZone)
	endTime, validEnd := parseTime(bhs.Spec.EndTime, bhs.Spec.TimeZone)
	if !validStart || !validEnd {
		logger.Error(errors.New("invalid time format"), "StartTime or EndTime has invalid format or time zone, expected HH:mm:ss")
		return ctrl.Result{}, errors.New("invalid StartTime or EndTime format, expected HH:mm:ss")
	}
	startTime = normalizeToAnchorDate(startTime)
	endTime = normalizeToAnchorDate(endTime)

	// Get the current time with time zone adjustments
	currentTime, err := r.determineCurrentTime(bhs)
	if err != nil {
		logger.Error(err, "Failed to determine current time")
		return ctrl.Result{}, err
	}
	currentTime = normalizeToAnchorDate(currentTime)

	// Determine business hours, including cross-day cases
	isBusinessHours := isWithinBusinessHours(currentTime, startTime, endTime)
	reason := "OutsideBusinessHours"
	desiredReplicas := bhs.Spec.MinReplicas
	if isBusinessHours {
		reason = "WithinBusinessHours"
		desiredReplicas = bhs.Spec.MaxReplicas
	}

	// List Deployments and apply scaling
	selector, err := metav1.LabelSelectorAsSelector(&bhs.Spec.DeploymentSelector)
	if err != nil {
		logger.Error(err, "Invalid label selector")
		return ctrl.Result{}, err
	}
	deploymentList := &appsv1.DeploymentList{}
	if err := r.List(ctx, deploymentList, &client.ListOptions{Namespace: bhs.Namespace, LabelSelector: selector}); err != nil {
		logger.Error(err, "Failed to list Deployments")
		return ctrl.Result{}, err
	}

	for _, deployment := range deploymentList.Items {
		if *deployment.Spec.Replicas != desiredReplicas {
			oldCount := *deployment.Spec.Replicas
			deployment.Spec.Replicas = &desiredReplicas
			if err := r.Update(ctx, &deployment); err != nil {
				logger.Error(err, "Failed to update Deployment replicas", "Deployment", deployment.Name)
				r.Recorder.Eventf(bhs, corev1.EventTypeWarning, ScaleFailure, "Failed to scale Deployment %s from %d to %d replicas: %v", deployment.Name, oldCount, desiredReplicas, err)
				continue
			}
			r.Recorder.Eventf(bhs, corev1.EventTypeNormal, ScaleSuccess, "Scaled Deployment %s from %d to %d replicas. Reason: %s", deployment.Name, oldCount, desiredReplicas, reason)
			logger.Info("Updated Deployment replicas", "Deployment", deployment.Name, "Replicas", desiredReplicas)
		}
	}

	// Update status
	if bhs.Status.CurrentReplicas != desiredReplicas || bhs.Status.LastUpdated.IsZero() {
		bhs.Status.CurrentReplicas = desiredReplicas
		bhs.Status.LastUpdated = metav1.Now()
		if err := r.Status().Update(ctx, bhs); err != nil {
			logger.Error(err, "Failed to update BusinessHoursScaler status")
			return ctrl.Result{}, err
		}
	}

	logger.Info("Reconciliation complete", "TimeUntilNextReconcile", intervalBetweenReconciles)
	return ctrl.Result{RequeueAfter: intervalBetweenReconciles}, nil
}

func (r *BusinessHoursScalerReconciler) determineCurrentTime(bhs *appsv1alpha1.BusinessHoursScaler) (time.Time, error) {
	loc := time.Local

	timeZone := bhs.Spec.TimeZone
	l, err := time.LoadLocation(timeZone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time zone '%s': %w", timeZone, err)
	}
	loc = l

	if fakeTimeStr, exists := bhs.Annotations[fakeTimeAnnotationKey]; exists {
		// Use fake time with specified time zone
		fakeTime, err := time.ParseInLocation(time.TimeOnly, fakeTimeStr, loc)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid fake-time format '%s', expected HH:mm:ss", fakeTimeStr)
		}
		return fakeTime, nil
	}

	// Use real clock with specified time zone
	return time.Now().In(loc), nil
}

func parseTime(timeStr string, timeZone string) (time.Time, bool) {
	zeroTime := time.Time{}

	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return zeroTime, false
	}
	t, err := time.ParseInLocation(time.TimeOnly, timeStr, loc)
	if err != nil {
		return zeroTime, false
	}
	return t, true
}

// Used for removing all year, month and date from a given time.Time. This is useful for comparing time.Time instances.
func normalizeToAnchorDate(t time.Time) time.Time {
	anchorDate := time.Date(2000, 1, 1, t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	return anchorDate
}

func isWithinBusinessHours(currentTime, startTime, endTime time.Time) bool {
	if startTime.After(endTime) {
		// Cross-day: business hours span midnight
		return currentTime.After(startTime) || currentTime.Before(endTime)
	}
	// Regular: business hours within the same day
	return currentTime.After(startTime) && currentTime.Before(endTime)
}
