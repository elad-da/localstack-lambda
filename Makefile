.PHONY: deps run clean

deps: npm install --save-dev concurrently nodemon serverless serverless-localstack
	docker pull localstack/localstack:1.0.2	

run: 
	-docker-compose up -d && sleep 20
	-npm start

clean:
	-docker-compose down -v
	-rm -rf bin