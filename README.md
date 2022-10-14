# Wanted Vehicles

## Development

Build the binary

```sh
make
```

Start postgres

```sh
docker-compose up -Vd postgres
```

Run sql migrations

```sh
migrate -source file://migrations -database postgres://postgres:password@127.0.0.1/wanted\?sslmode=disable up
```

Run the web server

```sh
./bin/server
```

## Test

Start postgres

```sh
docker-compose up -Vd postgres
```

Run sql migrations

```sh
migrate -source file://migrations -database postgres://postgres:password@127.0.0.1/wanted\?sslmode=disable up
```

Run tests

```sh
go test -race -bench=. -v ./...
```

## Usage

For example, get information about this amazing Tesla Model S

```sh
http http://localhost:8080/api/v1/wanted/vehicles?number=СВ5501ВХ
```

```json
[
    {
        "body_number": "5YJSA1E28HF176944",
        "brand": "TESLA - MODEL S",
        "color": "СІРИЙ",
        "id": "3019228562749883",
        "insert_date": "2019-08-16T15:37:54Z",
        "kind": "ЛЕГКОВИЙ",
        "number": "СВ5501ВХ",
        "ovd": "СОЛОМ’ЯНСЬКЕ УПРАВЛІННЯ ПОЛІЦІЇ ГУНП В М. КИЄВІ",
        "revision_id": "17082019_1",
        "status": "removed",
        "theft_date": "2019-08-16"
    }
]
```

## License

Project released under the terms of the MIT [license](./LICENSE).
