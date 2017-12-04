#Grab the latest alpine image
FROM alpine:3.6

# Add server and client bin
RUN mkdir -p /dist
ADD ./dist/server /server
ADD ./dist/client /dist/client-latest

# Run the image as a non-root user
RUN adduser -D myuser
USER myuser

# Run the app.  CMD is required to run on Heroku
# $PORT is set by Heroku			
CMD /server