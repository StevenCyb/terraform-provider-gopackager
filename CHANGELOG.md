## 0.2.2
FEAT:
- Hashes for ZIP files are now computed by the content of the ZIP instead of the ZIP files itself

FIXES:
- Fix the SHA512 hash

CHORE:
* General improvements (can include including docs, code formatting, comments, etc.)

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
