package rule

import future.keywords.if
import future.keywords.contains

custom contains result if {
    val := custom_function("param", "another")
    #val := base64.decode("ZXhwZWN0ZWQ=")

    val == "expected"

    result = "ok"
}
