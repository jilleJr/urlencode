# urlencode

Super basic URL encoding utility. I needed one, so I decided to make one.

![Screenshot from 2021-10-05 18-34-51](https://user-images.githubusercontent.com/2477952/136065171-c4c693f4-38de-4b3b-b628-f066b4a96769.png)

## Installation

Requires Go v1.16 (or higher)

```console
$ go install github.com/jilleJr/urlencode
```

## Features

- Encodes

- Decodes

- Colored output to highlight what's encoded/decoded

- Read from STDIN or from a file

## Usage

```console
$ urlencode --help
urlencode v1.0.0  Copyright (C) 2021  Kalle Jillheden

  License GPLv3+: GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>.
  This program comes with ABSOLUTELY NO WARRANTY; for details type '--license-w'
  This is free software, and you are welcome to redistribute it
  under certain conditions; type '--license-c' for details.

Encodes/decodes the input value for HTTP URL by default and prints
the encoded/decoded value to STDOUT.

  urlencode             // read from STDIN
  urlencode myfile.txt  // read from myfile.txt

Flags:
  -d, --decode            decodes, instead of encodes
  -e, --encoding string   encode/decode format (default "path-segment")
  -h, --help              show this help text and exit
  -l, --lines             encode/decode each line by themselves
      --version           show version and exit

Valid encodings (--encoding):
 SHORT  LONG          EXAMPLE
                      http://user:pass@site.com/index.html?foo=bar#Hello
 s      path-segment  --------------------------index.html--------------
 p      path          -------------------------/index.html--------------
 q      query         -------------------------------------foo bar------
 h      host          -----------------site.com-------------------------
 c      cred          -------user:pass----------------------------------
 f      frag          --------------------------------------------#Hello

                      http://[::1%25eth0]/home/index.html
 z      zone          --------------eth0-----------------
```

## License

Copyright &copy; 2021 Kalle Jillheden

License GPLv3+: GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>.
See full license text in the [LICENSE](./LICENSE) file.
