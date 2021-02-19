# Build/Run Locally

### Download dependencies
Download the following CLIs for your OS and add to your path
- BOSH CLI (https://github.com/cloudfoundry/bosh-cli/releases)
- PKS CLI (https://network.pivotal.io/products/pivotal-container-service)
- GOVC CLI (https://github.com/vmware/govmomi/releases)

### Build the app
Build using newer GO (built with go1.15+)
- `go build -o platform-info`
- `chmod +x platform-info`

# Build Docker Image
Download the following Linux flavors of the CLIs and add to the clis directory 
- BOSH CLI (https://github.com/cloudfoundry/bosh-cli/releases)
- PKS CLI (https://network.pivotal.io/products/pivotal-container-service)
- GOVC CLI (https://github.com/vmware/govmomi/releases)

To build the Docker image, run `docker build -t <your repo>/platform-info .`

Copy the `.platform-info-config.yaml` file to your home directory and fill out the values

Example for running the image using the config file

`docker run -v <path to your OpsMan CA Cert>:/mnt/cert -v ~/.platform-info-config.yaml:/mnt/.platform-info-config.yaml <your docker repo>/platform-info:latest /usr/bin/platform-info tkgi --config /mnt/.platform-info-config.yaml`

Check the help for more options `./platform-info --help`

If you're feeling adventurous, you can also provide all the credential options as environment variables or flags to the CLI 
