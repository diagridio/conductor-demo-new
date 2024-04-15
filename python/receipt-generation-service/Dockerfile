FROM --platform=$BUILDPLATFORM python:3.8-alpine

WORKDIR /src

COPY requirements.txt .
RUN pip3 install -r requirements.txt

COPY . .

EXPOSE 80


CMD [ "python", "./main.py" ]