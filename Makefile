.PHONY: run
run:
	kubectl kustomize local/ --enable-helm | kubectl apply --server-side --force-conflicts -f -