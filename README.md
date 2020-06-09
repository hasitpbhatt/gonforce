# gonforce
Gonforce ( Read as Go Enforce ) enforces the given package follows the guidelines defined in "gonforce.yaml" file.

# Install

```
go get -u github.com/hasitpbhatt/gonforce
```
Once installed, use gonforce command in the package containing "gonforce.yaml"

# Samples
Sample files may look something like the following

1. Default Blocklist
```
package: github.com/hasitpbhatt/gonforce

default:
  type: blocklist
  imports:
    - fmt
    - gopkg.in
  except:
    - gopkg.in/yaml.v2
```

Here we are ensuring the package can't use any imports starting with "gopkg.in" or "fmt", but with an exception "gopkg.in/yaml.v2" which would not be allowed otherwise due to "gopkg.in" being one of the blocklisted imports.

2. Default Allowlist
```
package: github.com/hasitpbhatt/gonforce

default:
  type: allowlist
  imports:
    - fmt
    - gopkg.in
  except:
    - gopkg.in/models
```

Here we are ensuring only imports that can be used should always start with "fmt" or "gopkg.in", but we don't want to allow "gopkg.in/models" due to some specific reason. Any other imports will automatically fail the gonforce validation.

3. Basic Model-Controller Configuration
```
package: github.com/hasitpbhatt/gonforce

default:
  type: allowlist
  imports:
    - fmt
    - gopkg.in
    - github.com/hasitpbhatt/gonforce
  except:
    - gopkg.in/models

rules:
  - name: models
    rule:
      type: blocklist
      imports:
        - github.com/hasitpbhatt/gonforce/controllers # models can't import controller package

  # Even if this was uncommented, it would serve the same purpose for there is no restrictions 
  # on controllers at this moment
  # - name: controllers
  #   rule:
  #     type: blocklist
```

Here, along with ensuring the we only use specific packages inside the default, using blocklists we ensure that we don't endup calling an subpackage of controller inside models package. And controllers doesn't have any extra restrictions other than already imposed by default package rule.
