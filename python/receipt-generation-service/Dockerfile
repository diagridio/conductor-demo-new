FROM --platform=$BUILDPLATFORM python:3.11

WORKDIR /src

COPY requirements.txt .
RUN pip3 install --no-cache-dir --upgrade -r requirements.txt

COPY . .

EXPOSE 80

CMD [ "python", "./app/main.py" ]0
