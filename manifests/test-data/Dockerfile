FROM mongo
WORKDIR /
RUN apt update && \
    apt -y install wget unzip && \
    wget https://github.com/pasiol/remainders-user/releases/download/v0.1.1/remainders-user && \
    chmod u+x ./remainders-user
COPY test_data.json.zip /
CMD [""]