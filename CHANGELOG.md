<a name="v0.5.0"></a>
# [v0.5.0](https://github.com/qri-io/starlib/compare/v0.4.2...v0.5.0) (2021-04-29)

version 0.5.0 breaks the APIs of `math`, `time`, and `json`.

After a *large* community effort, we've successfully moved two packages from starlib into go-starlark:
* math: google/starlark-go#357
* time: google/starlark-go#327

Both of these packages have had a great deal of vetting & improvement in the process of migrating. While this does introduce many breaking changes for both packages, both packages have better stability & performance characteristics.

Moving forward, we'll try to follow the pattern of developing & testing packages here, and for those deemed worthy, move them up into go-starlark. Once a package lands there, we'll switch starlib to being a strict import-and-document-only, which effectively locks their API. **Keeping package APIs matched to go-starlark will cut down on drift, benefitting the broader starlark ecosystem**. We're currently working on doing the same with the `re` / regexp package.

In the same vein, we've switched to the go-starlark `json` module, which adds another breaking API change for this release.

### Bug Fixes

* **http:** send multipart form data correctly ([c365bf6](https://github.com/qri-io/starlib/commit/c365bf6))
* **json:** Fix unmarshaling of 64 bit integers ([#71](https://github.com/qri-io/starlib/issues/71)) ([86d6e56](https://github.com/qri-io/starlib/commit/86d6e56))
* **math,time:** restore ".star" suffixes to math & time modules ([4528225](https://github.com/qri-io/starlib/commit/4528225))
* **time:** fix location parsing in time() ([aa672aa](https://github.com/qri-io/starlib/commit/aa672aa))
* **util:** handle tuples correctly ([#61](https://github.com/qri-io/starlib/issues/61)) ([863ced0](https://github.com/qri-io/starlib/commit/863ced0))
* **version:** correct erroneous version indicator ([#66](https://github.com/qri-io/starlib/issues/66)) ([ce3581b](https://github.com/qri-io/starlib/commit/ce3581b))


### Code Refactoring

* **math,time:** replace math & time modules with Starlark upstream versions ([c7a7f60](https://github.com/qri-io/starlib/commit/c7a7f60))


### Features

* **http:** support URL encoding for forms ([5c99138](https://github.com/qri-io/starlib/commit/5c99138))
* **json:** replace encoding/json with go.starlark.net/starlarkjson ([f21f027](https://github.com/qri-io/starlib/commit/f21f027))
* **time:** Add strftime method to time instances ([08db895](https://github.com/qri-io/starlib/commit/08db895))
* **time:** make string representation reversible ([5d9d07b](https://github.com/qri-io/starlib/commit/5d9d07b))
* **util:** unmarshal/marshal for starlib/time.Time to go's standard time.Time and back ([1364b0c](https://github.com/qri-io/starlib/commit/1364b0c))


### BREAKING CHANGES

* **json:** overhaul json module. See package doc for details
* **math,time:** math & time modules have been overhauled. Refer to package documentation for
details



<a name="v0.4.2">v0.4.2</a>
# [v0.4.2](https://github.com/qri-io/starlib/compare/v0.4.1...v0.4.2) (2020-06-29)

This release brings a number of enhancements to the `time`, `re`, and golang-side utility packages. It also introduces two new packages: `encoding/yaml` and `hash`. Here's the docs for the two new packages:

### yaml
yaml provides functions for working with yaml data

#### Functions

##### `dumps(obj) string`
serialize obj to a yaml string

**parameters:**

| name | type | description |
|------|------|-------------|
| `obj` | `object` | input object |


##### `loads(source) object`
read a source yaml string to a starlark object

**parameters:**

| name | type | description |
|------|------|-------------|
| `source` | `string` | input string of yaml data |


### hash
hash defines hash primitives for starlark.

#### Functions

##### `md5(string) string`
returns an md5 hash for a string

##### `sha1(string) string`
returns an sha1 hash for a string

##### `sha256(string) string`
returns an sha256 hash for a string


### Bug Fixes

* **time:** unix() and unix_nano() returns 0 for epoch ([c50ebc2](https://github.com/qri-io/starlib/commit/c50ebc2))


### Features

* **customType:** Added support for custom type and fixed bugs ([#49](https://github.com/qri-io/starlib/issues/49)) ([c32c667](https://github.com/qri-io/starlib/commit/c32c667)), closes [#48](https://github.com/qri-io/starlib/issues/48)
* **encoding/yaml:** encoding yaml package based on gopkg.in/yaml.v2 ([68e22bc](https://github.com/qri-io/starlib/commit/68e22bc))
* **hash:** add hash module ([686ae7b](https://github.com/qri-io/starlib/commit/686ae7b))
* **re:** add compile function to regex package ([6fe15cd](https://github.com/qri-io/starlib/commit/6fe15cd))
* **starlib:** util.Marshal extended compatibility and tests ([a310f83](https://github.com/qri-io/starlib/commit/a310f83))
* **time:** Add fromtimestamp method ([4e7be49](https://github.com/qri-io/starlib/commit/4e7be49))
* **time:** Add in_location and format methods ([316e7aa](https://github.com/qri-io/starlib/commit/316e7aa))
* **time:** Add methods hours(), minutes(), etc to duration ([f374e23](https://github.com/qri-io/starlib/commit/f374e23))
* **time:** Add unix and unix_nano methods ([7c32cb7](https://github.com/qri-io/starlib/commit/7c32cb7))



<a name="v0.4.1"></a>
# [v0.4.1](https://github.com/qri-io/starlib/compare/v0.4.0...v0.4.1) (2019-06-10)

Quick patch release the brings the `get_text` function to the bsoup package.

### Bug Fixes

* **version:** correct erroneous version indicator ([709ae5e](https://github.com/qri-io/starlib/commit/709ae5e))


### Features

* **bsoup:** add get_text method to bsoup ([#31](https://github.com/qri-io/starlib/issues/31)) ([499d3e8](https://github.com/qri-io/starlib/commit/499d3e8))



<a name="0.4.0"></a>
# [0.4.0](https://github.com/qri-io/starlib/compare/v0.3.2...v0.4.0) (2019-05-29)

In preparation for go 1.13, in which `go.mod` files and go modules are the primary way to handle go dependencies, we are going to do an official release of all our modules. This will be version v0.4.0 of `starlib`. That said, this release also includes a new package for working with HTML:

## new beautiful soup-like HTML package
Our `html` package is difficult to use, and we plan to deprecate it in a future release. In it's place we've introduced `bsoup`, a new package that implements parts of the [beautiful soup 4 api](https://www.crummy.com/software/BeautifulSoup/bs4). It's _much_ easier use, and will be familiar to anyone coming from the world of python.


### Bug Fixes

* **bsoup:** Multiple fixes, cleanups. Test-only code. Own TODOs. ([dd1c8ba](https://github.com/qri-io/starlib/commit/dd1c8ba))


### Features

* **bsoup:** Bsoup is a wrapper to imitate beautifulSoup in starlark ([3c0caeb](https://github.com/qri-io/starlib/commit/3c0caeb))



<a name="0.3.2"></a>
# [0.3.2](https://github.com/qri-io/starlib/compare/v0.3.1...v0.3.2) (2019-04-03)


### Bug Fixes

* **xlsx:** fix updated excelize api ([cfe3c52](https://github.com/qri-io/starlib/commit/cfe3c52))



<a name="0.3.1"></a>
# [0.3.1](https://github.com/qri-io/starlib/compare/v0.3.0...v0.3.1) (2019-03-14)


### Features

* **re.search,csv.write_all:** add csv.write_all and re.search functions ([3eccb40](https://github.com/qri-io/starlib/commit/3eccb40))



<a name="0.3.0"></a>
# [0.3.0](https://github.com/qri-io/starlib/compare/v0.2.0...v0.3.0) (2019-03-07)

Added a JSON package, fixed math not being accessible in the default loader, a number of minor refactors


### Bug Fixes

* **math:** add missing math module to loader, add math.round(x) ([7923d11](https://github.com/qri-io/starlib/commit/7923d11))


### Features

* **json:** add initial json-package ([c165930](https://github.com/qri-io/starlib/commit/c165930))



<a name="0.2.0"></a>
# 0.2.0 (2019-01-22)

This is the first proper release of starlib. Packages added this release:
* encoding/base64.star
* encoding/csv.star
* geo.star
* http.star
* math.star
* re.star
* time.star
* xlsx.star
* zipfile.star

### Bug Fixes

* **time:** fix time errors ([65d5dd3](https://github.com/qri-io/starlib/commit/65d5dd3))


### Code Refactoring

* **.star:** change userland file endings to .star ([daa23e1](https://github.com/qri-io/starlib/commit/daa23e1))


### Features

* **binary ops:** initial binary operators for time & duration ([9e1f4d3](https://github.com/qri-io/starlib/commit/9e1f4d3))
* **csv:** add initial single-function csv package ([0532a9f](https://github.com/qri-io/starlib/commit/0532a9f))
* **csv params:** add arguments to configure csv.read_all ([61c74d6](https://github.com/qri-io/starlib/commit/61c74d6))
* **encoding/base64:** add basic base64 encode/decode module ([4204e76](https://github.com/qri-io/starlib/commit/4204e76))
* **geo.MultiPolygon:** need MultiPolygon ([eba9b24](https://github.com/qri-io/starlib/commit/eba9b24))
* **geo.parseGeoJSON:** add basic support for paring geoJSON strings ([9368420](https://github.com/qri-io/starlib/commit/9368420))
* **geo.Point:** implement basic geo.point type ([a7b7bf2](https://github.com/qri-io/starlib/commit/a7b7bf2))
* **html,http,xlsx:** add new stub packages ([d3f9b53](https://github.com/qri-io/starlib/commit/d3f9b53))
* **http:** accept better http params for forming requests ([c97ba10](https://github.com/qri-io/starlib/commit/c97ba10))
* **math:** add initial math package ([e41fca9](https://github.com/qri-io/starlib/commit/e41fca9))
* **math:** Mathematical constants and comments ([55037eb](https://github.com/qri-io/starlib/commit/55037eb))
* **re:** initial regexp library ([3c90541](https://github.com/qri-io/starlib/commit/3c90541))
* **starlib:** make this repo a stub for our standard library ([2317ac5](https://github.com/qri-io/starlib/commit/2317ac5))
* **time:** initial work on time module ([f5c91a8](https://github.com/qri-io/starlib/commit/f5c91a8))
* **time.sleep:** add sleep method ([2aae5ff](https://github.com/qri-io/starlib/commit/2aae5ff))
* **zipfile:** initial zipfile package ([ec3ff67](https://github.com/qri-io/starlib/commit/ec3ff67))


### BREAKING CHANGES

* **.star:** all imports now end with '.star' instead of '.sky'



