#
# @Author: sejal chougule 
# @Date: 2018-05-02 17:21:41 
# @Last Modified by:   sejal chougule 
# @Last Modified time: 2018-05-02 17:21:41 
# 

VENDOR_DIR="/root/deploynow-dependencies/vendor"

run: clean
	ENV=development fresh -c dnow_runner.conf

run-dev: clean
	ENV=development go run main.go

run-prod: clean
	ENV=production go run main.go

build: clean
	ENV=developement GOOS=linux GOARCH=amd64 go build .

vendor: clean
	@echo "\n -> Copying vendor dependencies -> started" 
	cp -R ${VENDOR_DIR} ${PWD}
	@echo "\n -> Copying vendor dependencies -> finished" 

# dev: clean
# 	@echo "\n -> Copying vendor dependencies -> started" 
# 	cp -R /Users/sejal/Projects/Personal/go/src/github.com/sjljrvis/deploynow/vendor /Users/sejal/Projects/Personal/go/src/github.com/sjljrvis/deploynow-dependencies
# 	@echo "\n -> Copying vendor dependencies -> finished" 

clean:
	@echo "\n -> Cleaning cache and log files\n" 
	-find . -name 'nohup.out' -delete
	@echo "\n -> Cleaning done\n"

help:
	@echo "\nPlease call with one of these targets:\n"
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F:\
        '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}'\
        | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs | tr ' ' '\n' | awk\
        '{print "    - "$$0}'
	@echo "\n"	
