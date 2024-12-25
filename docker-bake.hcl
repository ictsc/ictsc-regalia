group "default" {
  targets = ["backend"]
}

target "docker-metadata-action" {}

target "backend" {
  inherits = ["docker-metadata-action"]
  context = "./backend"
  matrix = {
    image = [
      "scoreserver",
      "toolbox"
    ]
  }
  name = "backend-${image}"
  tags = make_tags("${image}")
  target = "${image}"
}

variable "DOCKER_METADATA_OUTPUT_TAGS" {
    default = ""
}
function "make_tags" {
    params = [ns]
    result = split("\n", replace("${DOCKER_METADATA_OUTPUT_TAGS}", ":", "/${ns}:"))
}
