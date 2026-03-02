#' Hello World
#'
#' A simple example function.
#'
#' @param name A character string with a name to greet.
#' @return A character string with a greeting.
#' @export
#' @examples
#' hello("World")
hello <- function(name = "World") {
    paste0("Hello, ", name, "!")
}
