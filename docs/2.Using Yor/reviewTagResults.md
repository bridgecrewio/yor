# Reviewing Tag Results

```sh
./yor tag -d . -o cli
# default cli output

./yor tag -d . -o json
# json output

./yor tag -d . --output cli --output-json-file result.json
# will print cli output and additional output to file on json file -- enables programatic analysis alongside printing human readable result
```