group "default" {
  targets = ["backend", "frontend"]
}

target "docker-metadata-action" {
  tags = []
}

target "base" {
  args = target.docker-metadata-action.args
  labels = target.docker-metadata-action.labels
  annotations = target.docker-metadata-action.annotations
}

target "backend" {
  inherits = ["base"]
  context = "./backend"
  matrix = {
    image = [
      "scoreserver",
      "toolbox"
    ]
  }
  name = "backend-${image}"
  tags = make_tags(target.docker-metadata-action.tags, "${image}")
  target = "${image}"
}

target "frontend" {
  inherits = ["base"]
  context = "./frontend"
  matrix = {
    image = ["frontend"]
  }
  name = "frontend-${image}"
  tags = make_tags(target.docker-metadata-action.tags, "${image}")
  target = "${image}"
}

function "make_tags" {
    params = [tags, name]
    result = split("\n", replace(join("\n", tags), ":", "/${name}:"))
}
