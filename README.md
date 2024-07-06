# Card Store
Storefront for querying, selling, buying and trading collectable cards.

## Projects

### server
Backend service, stores user data, card data and collection data in a PostgreSQL database, uses Valkey for caching card and query data.

Requires running Docker containers for the database, cache and query cache.

Run the backend using the commands
```sh
cd server
go run .
```

### client
Website frontend, uses React + Bootstrap 5.

Run the frontend using the commands
```sh
cd client
npm run dev
```

### scripts
Scripts used to test/help in the development of the marketplace
* _pop-dv.sh_ - takes a _.sql_ file as an argument, copies that file to the database container, then runs it
* _generate-mtg-sql.py_ - python script for generating a _.sql_ file that contains cards from the collectable card game Magic: the Gathering.
