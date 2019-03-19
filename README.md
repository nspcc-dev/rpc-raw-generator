# RPC raw generator

This is application for generating HTTP-requests flow of the prepared raw transactions. 
It may be used with neo-bench environment. For more details read here.

## Usage
Put your file with raw transactions to the `raw` directory and start application with `make up`

```
$ cp ~/raw.txs ./raw
$ make up
. . .
enter inside container:

/var/test # rpc-generator --help
Usage of rpc-generator:
  -c int
        number of simultaneous connections (default 1)
  -f string
        file with raw transactions (default "./raw.txs")
  -t int
        total amount of transactions (default 100)
  -url string
        connection url (default "http://127.0.0.1:30334/")
/var/test # rpc-generator -c 1 -f ./raw.txs -t 100 -url "http://127.0.0.1:30334/"
2019/03/19 07:16:38 Sending 100 transactions (100 tx in 1 connections) to http://127.0.0.1:30334/
2019/03/19 07:16:38 connection  1 ) time:  245.639094ms
2019/03/19 07:16:38 ---
2019/03/19 07:16:38 Total time:  245.678656ms  Approximate  TPS:  407.0357662653446
2019/03/19 07:16:38 For accurate TPS values use external connection sniffers, e.g. Wireshark 
/var/test # ^C
```

## License

This project is licensed under the GPL v3.0 License - see the 
[LICENSE](LICENSE) file for details
