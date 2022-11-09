.PHONY: run
run:
	kubectl kustomize local/ --enable-helm | kubectl apply -f -