go test -v -cover -coverprofile=cover -coverpkg=../container,../convert,../dispose,../errs,../odi,../resolve,../types
go tool cover -html=cover -o cover.html