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



