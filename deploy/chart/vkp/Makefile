bin/helm-docs:
	mkdir -p bin/
	wget -qO- https://github.com/norwoodj/helm-docs/releases/download/v1.11.0/helm-docs_1.11.0_Linux_x86_64.tar.gz | tar -xvz -C bin/ helm-docs
	chmod +x bin/helm-docs

.PHONY: docs
docs:
	bin/helm-docs