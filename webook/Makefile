.PHONY: docker
docker:
	@rm webook || true # 先判断是否存在.\webook 可执行文件
	@go build -tags=k8s -o webook . # 生成 webook 可执行文件
	@docker rmi -f flycash/webook:v0.0.1 # 删除原来的版本，前提是有代码有改动
	@docker build -t flycash/webook:v0.0.1 . # 生成 v0.0.1 版本