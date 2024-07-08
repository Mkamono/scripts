# ozo-control

## install

### apple silicon Mac

```shell
sudo curl -o /usr/local/bin/ozo-control https://raw.githubusercontent.com/Mkamono/scripts/main/ozo-control/assets/darwin-arm64/ozo-control && \
sudo chmod +x /usr/local/bin/ozo-control
```

uninstall

```shell
ozo-control clean && sudo rm /usr/local/bin/ozo-control
```

## Usage

```shell
ozo-control --help
```

### Check In

```shell
ozo-control i
```

or

```shell
ozo-control check-in
```

### Check Out

```shell
ozo-control o
```

or

```shell
ozo-control check-out
```

### Register Holiday

As default, it registers weekends and public holidays in Japan.

```shell
ozo-control r
```

or

```shell
ozo-control register-holiday
```

### Override Option

Override option are available. To override already registered holidays, use `-o` option.

```shell
ozo-control r -o
```

or

```shell
ozo-control register-holiday -o
```
