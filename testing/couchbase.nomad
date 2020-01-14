job "couchbase" {
  datacenters = ["dc1"]

  group "couchbase" {
    count = 2

    task "couchbase" {
      driver = "docker"

      config {
        image = "couchbase"

        port_map {
          port0 = 11210
          port1 = 8091
          port2 = 8092
          port3 = 8093
          port4 = 8094

        }

      }

      resources {
        network {
          mbits = 10
        }
      }

      // service {
      //   name = "couchbase"
      //   port = "port1"
      // }
    }
  }
}