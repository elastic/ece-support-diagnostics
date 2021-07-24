## How to prepare a new release : 

1- Update the version in `src/ece-diagnostics.sh` and commit if not done on line 3: 
```
ECE_DIAG_VERSION=2.0.0
```

2- Create a zip file for the release with version : 
```
version=2.0.0
cd src
zip "$HOME/ece-support-diagnostics-v${version}.zip" ece-diagnostics.sh
```

3- Create new release by uploading the zip file and using same version tag


## testing

No automated testing at the moment, reference tests :
```
./ece-diagnostics.sh -s -sp /invalid/path #test path is corrected automatically (folder elatic will contain data - if fails then execution will stop)

./ece-diagnostics.sh -u admin -p wrongpassword #execution should stop with wrong password

./ece-diagnostics.sh -u admin # should prompt for password

./ece-diagnostics.sh -u admin -p tmk7w77tsfLB2sGijhn8SHnGr6emXfHEkklUwCCrS89 -de 9ac44b53573d4fedaa768ef0453e985a,a4c1cbccb5f3401a9627e39516113771 #test returning plan activity logs for 2 deployments (log will show some call including show_plan_logs=true&show_plan_history=true)

./ece-diagnostics.sh -u admin -p tmk7w77tsfLB2sGijhn8SHnGr6emXfHEkklUwCCrS89 -de 9ac44b53573d4fedaa768ef0453e985a,a4c1cbccb5f3401a9627e39516113771 -lh 1 -d -s #test all options, should only get today logs and last 1 hour of container logs, note some system commands may fail with a warning currently
```