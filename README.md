# cbr_currencies

`cbr_currencies` is a tool to get the Bank of Russia exchange rate for today or specified date.


## Usage

Linux:

```
./cbr_currencies
```

Windows:

```
cbr_currencies.exe
```

If you only need to get some currencies, specify the flag '-c' and then currency codes (according to ISO 4217) separated by commas:

```
./cbr_currencies -c usd,eur
```

If you need to get data for another dates, specify the flag '-d' and then dates (in format 'day.month.year') separated by commas:

```
./cbr_currencies -d 4.03.20,10.12.20
```

If you need to save data, specify the flag '-s' and then a name of the SQLite database file in which the exchange rate data should be saved:

```
./cbr_currencies -s currencies.db
```


## License

The code is under the MIT license.
