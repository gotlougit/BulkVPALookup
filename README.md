# BulkVPALookup

## Overview

This CLI tool (written in Go) helps look up a whole bunch of Indian cell phone numbers you may have obtained (eg from WhatsApp group membership lists), and returns a VCF full of names of people using those phone numbers, which you can import into your contacts list and hence communicate easily with others.

Rather than manually saving everyone's number, this tool automates that procedure.

## Usage

1. Clone this repo

```
git clone https://git.sr.ht/~gotlou/BulkVPALookup
```

2. cd into the newly created BulkVPALookup/ folder

```
cd BulkVPALookup/
```

3. Run the following command: 

```
go run main.go phonenums.txt contacts.VCF
```

where phonenums.txt is a file containing a list of all phone numbers you want to lookup in this format: 

```
9999999999
8888888888
7777777777
...
```

and contacts.vcf is the name of the VCF you want the results of the search to be saved.

## How does it work

It uses a free API to guess at people's UPI VPAs and gets the name if it exists. After getting most of the names this way, it writes what it has found into a VCF, which is a standard for saving contacts that almost all contacts apps will recognize and import.

To learn more about how this works, read [my blog post](https://gotlou.srht.site/phone-num-lookup.html) on this topic.

## Acknowledgements

Thanks to [Aseem Shrey](https://aseemshrey.in/) for building a similar tool [here](https://github.com/LuD1161/upi-recon-cli), written in Go. One of the GitHub issues on that page led me to upibankvalidator.com

## Disclaimer

I do NOT own or operate or have anything to do with upibankvalidator.com. While they don't really get a whole lot of info about you specifically other than that you made the request using this tool and what IP address you had at the time, I don't know how they use this information. Use at your own risk, I am NOT liable for any damages.

This tool was primarily made for educational purposes.

## License

This code is licensed under the GNU General Public License v2.
