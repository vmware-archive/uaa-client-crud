clean:
	rm -rf ./product
	rm -rf ./release
	rm ./resources/uaa-crud-release.tgz

build:
	wget https://github.com/cf-platform-eng/uaa-client-crud/releases/download/v0.7/uaa-crud-release.tgz -O ./resources/uaa-crud-release.tgz
	tile build
