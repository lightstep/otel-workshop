# otel-workshop-go
The base image for work on the OTel workshop. This will be used to demonstrate the basics and work through instrumentation steps. Slides can be found here:

https://docs.google.com/presentation/d/1m5NqQx_Z7R92Ri3poNq-QWSh54AB6M28zmQRH0xKmEg/edit?usp=sharing

Requirements:

This project currently requires Go 16 along with several packages to include:

go get go.opentelemetry.io/otel \                      
       go.opentelemetry.io/otel/trace 

go get go.opentelemetry.io/otel/sdk \                  
       go.opentelemetry.io/otel/exporters/stdout/stdouttrace

