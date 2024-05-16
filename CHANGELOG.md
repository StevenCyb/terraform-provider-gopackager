## 0.2.7
FEAT:
- New Git commit output hash is added where "*.go", "go.mod" or "go.sum" was changed.

## 0.2.6
CHORE:
- Update documentation for the new `provided.al2023` variant of lambda [see here](https://aws.amazon.com/blogs/compute/migrating-aws-lambda-functions-from-the-go1-x-runtime-to-the-custom-runtime-on-amazon-linux-2/).

FIXES:
- If destination is a directory that does not exists, it will be created. E.g. `/tmp/not_existing/binary` will create `/tmp/not_existing` directory.

## 0.2.5
CHORE:
- Base64 hashes are now computed based on the raw hash rather then on the hexadecimal encoded 

## 0.2.4
FIXES:
- ZIP hash is now computed based on the ZIP content and not contained files

## 0.2.3
CHORE:
- DOC improvement

## 0.2.2
FEAT:
- Hashes for ZIP files are now computed by the content of the ZIP instead of the ZIP files itself

FIXES:
- Fix the SHA512 hash

CHORE:
- General improvements (can include including docs, code formatting, comments, etc.)

## 0.2.1
FIXES:
- Working directory path on compiler is now defined correctly

## 0.2.0
FEAT:
- Output now has more hash variations

FIXES: 
- Handling working directory on compile correctly
- Fix for_each support by replacing `ValidateConfig` with `ConfigValidators`

## 0.1.1

FIXES:
- Compiler now uses working directory of the source
- Compiler now automatically installs mod packages
- Some small fixes

## 0.1.0

FEATURES:
- base functionality to provide compile feature
