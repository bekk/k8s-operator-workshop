# Case: Business Hours Scaler™
> ℹ️ Sørg for å ha installert nødvendig programvare som spesifisert [her](./readme.md) før du begynner.

## Bakgrunn
I denne workshopen skal du lage en Kubernetes-[operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) som inneholder et API for å skalere workloads basert åpningstider. Dette kan være nyttig for å spare kostnader ved å redusere antall pods i perioder hvor det er lite trafikk. Målet med denne workshopen er å gi deg en forståelse av hvordan en operator fungerer, og hvordan du kan lage en enkel operator selv.

Merk at det finnes allerede ulike måter å skalere instanser av applikasjoner på, f.eks. ved å skalere etter CPU/minne ([HPA](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)/[VPA](https://github.com/kubernetes/autoscaler/tree/master/vertical-pod-autoscaler)), basert på eventer med [KEDA](https://keda.sh/), eller ved å bruke egne metrikker. 

### Mappestruktur

Dette er et skjelett for operatoren, basert på Operator SDK (og kubebuilder) samt egne erfaringer. 

```
.
├── .dockerignore
├── .gitignore
├── .golangci.yml                                         # (Valgfri) anbefalt linting
├── .run                                                  # JetBrains run configurations
├── .vscode                                               # Visual Studio Code run configurations
├── Dockerfile
├── LICENSE
├── Makefile                                              # Ulike prekonfigurerte "common" tasks
├── api                                                   # Her defineres din CRD (API) mot konsumenter
│   └── v1alpha1                                          # API-versjon
│       ├── businesshoursscaler_types.go                  # Typer for API'et ditt
│       ├── groupversion_info.go                          # Registrering av v1alpha1 inn under en API group ("apps.k8s.bekk.no")
│       └── zz_generated.deepcopy.go                      # Autogenerert
├── bin                                                   # Diverse tooling, samt kompilert versjon av operatoren
├── cmd                                                   # Konvensjon: putt det man ønsker skal være en standalone executable som en mappe her
│   └── bekk-ws-operator                                  # Navnet på din operator
│       └── main.go                                       # Entrypoint
├── codegen.go                                            # Flagging til kompilatoren/go tooling vedr. kodegenerering
├── config                                                # Autogenererte ressurser (definisjoner)
│   ├── crd
│   │   └── apps.k8s.bekk.no_businesshoursscalers.yaml    # [autogen] Deploybar og delbar CRD. Kan brukes til lokal validering i editor.
│   └── rbac
│       └── role.yaml                                     # [autogen] En rolle som inneholder det som trengs for at operatoren skal fungere.
├── go.mod                                                # Go-dependencies. Som package.json eller pom.xml.
├── go.sum                                                # De faktiske dependenciene i bruk. Som package-lock.json.
├── internal                                              # En pakke som aldri vil bli eksponert for andre prosjekter (etter navnekonvensjon)
│   └── controller                                        # Ulike controllere (n=1 for vår del)
│       └── businesshoursscaler_controller.go             # Forretningslogikken din
├── readme.md
├── sample                                                # Ting å teste manuelt med
│   ├── 01_ns.yaml
│   ├── 02_deployment.yaml
│   └── 03_bhs.yaml
├── tests                                                 # Deklarative manifest-baserte tester
│   ├── config.yaml                                       # Config for Chainsaw
│   ├── [tester i egne mapper]
├── tools.go                                              # Tools for codegen/testing
└── workshop.md                                           # Denne filen
```
 
## Let's go 🏃‍♂️

### Steg 1a: The basics
1. Start med å klone repoet og åpne i din favoritt-editor/IDE.
   ```shell
   git clone https://github.com/bekk/k8s-operator-workshop.git
   ```
2. Sett opp et lokalt Kubernetes-cluster for lokal utvikling
   ```shell
   # Bruk kind til å lage et lokalt cluster (som Docker-containere på din maskin)
   make setup-local
   # Bytt til det nye clusteret
   kubectl config use-context kind-workshop --namespace=bekk-ws-operator-system
   ```

### Steg 1b: Kjøring

### Kjøring lokalt
Benytt en av de ferdige konfigurasjonene for enten VSCode eller IntelliJ/GoLand. Hvis ikke kan du kjøre:

```shell
# Generer CRD/RBAC-regler og putt de inn i clusteret, før operatoren kjøres opp lokalt
make run-local
```

### Kjøre tester
```shell
# Kjør tester (uten å kjøre opp operatoren, slik at du kan debugge)
make test

# Kjør tester (og operatoren)
make run-test

# Kjøre enkelttester (uten operator)
make test-single dir=tests/some-example-test

# Kjøre enkelttester (med operator)
make run-test TEST_DIR=tests/some-example-test
```

### Nullstill clusteret
Har du lyst på et cleant cluster?
```shell
kind delete cluster --name workshop && make setup-local
```

### Steg 2: Definer API-et ditt 📜
I mappen `api/v1alpha1` finner du `businesshoursscaler_types.go` some inneholder API'et du eksponerer til andre.

Definer felter for når applikasjonen skal skalere opp og når den skal skalere ned. Det er også å nødvendig å vite hvilken applikasjon (`Deployment`-ressurs) som skal skaleres. Hva annet trenger du? 🤔 

**Tips:** Se på annotasjoner for f.eks. [validering](https://book.kubebuilder.io/reference/markers/crd-validation) og dokumentasjonen for [kubebuilder](https://book.kubebuilder.io/cronjob-tutorial/new-api).

### Steg 3: Skriv forretningslogikken 🧠

Åpne controlleren din i `internal/controller/businesshoursscaler_controller.go`. Her er det funksjonen `Reconcile(ctx context.Context, req ctrl.Request)` som skal fylles ut.

Ting å tenke på når det kommer til implementasjon:
- Finnes `BusinessHoursScaler` objektet ditt når controlleren din får requesten?
- Kubernetes, Go og tid kan være en spennende kombinasjon
- [RBAC](https://book.kubebuilder.io/reference/markers/rbac) for andre ressurser

Ressurser:
- [Kubernetes Go client](https://pkg.go.dev/k8s.io/client-go)
- [Go: time](https://pkg.go.dev/time)

### Steg 4: Events 📢

En vanlig måte å kommunisere tilstand på i Kubernetes er ved å sende [events](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/event-v1/). Dette kan være nyttig for å gi informasjon til brukere, eller brukes i feilsøkingsøyemed.

Bruk en [`EventRecorder`](https://github.com/kubernetes/client-go/blob/master/tools/record/event.go#L93) til å si ifra om at man har skalert et `Deployment` og når det ikke går bra (feilscenarier). Hva er viktig å få med for å gi verdi? 

### Steg 5: Skriv tester 🧪

Man kan skrive tester for alt, og på forskjellige nivåer. I denne workshopen skal vi skrive deklarative tester på manifest-nivå. Dvs. at man kjører opp operatoren, applyer et manifest og asserter på at nye ressuer har kommet til/eksisterende har endret seg etc. Til dette bruker vi [Chainsaw](https://kyverno.github.io/chainsaw/latest/).

Se `tests/` for konfigurasjon og eksempeltest.

> 💡Tenk spesielt på hvordan du skal manipulere tid i testene.

### Steg 6: Metrikker 📊

Man kan få mange metrikker ut av boksen. Finn ut av hvordan man kan eksponere en egen HTTP-port med Prometheus-metrikker.

### Stretch goals

- Ta høyde for tidssoner (hvis du ikke har måttet gjøre dette tidligere)
- Definer en egen metrikk og observer at den blir eksponert
- Se om du kan få med deg hvis noen manuelt endrer på et `Deployment`, slik at du kan overstyre automatisk etter åpningstid. 
