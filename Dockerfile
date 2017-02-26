FROM iron/base
WORKDIR /app
EXPOSE 8080
COPY . /app/
CMD ["./drafthouse-seat-finder", "-baseUrl", "drafthouse"]
