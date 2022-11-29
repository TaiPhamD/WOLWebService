set(TARGET wolwebservice)

set(SRCS server.go)
SET(CC clang-cl)
SET(CXX clang-cl)

set(OUTPUT "${CMAKE_BINARY_DIR}/dist/wolwebservice")

# Custom go install 3rd party "Service" lib and build command
add_custom_command(OUTPUT golang_wolwebservice_build.out
        DEPENDS ${SRCS}
        WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
        COMMAND go mod download
        COMMAND ${CMAKE_COMMAND} -E echo "Compiling Golang wolwebservice project"
        COMMAND go build -o ${OUTPUT}   
)

add_custom_target(wolwebservice ALL
        ${CMAKE_COMMAND} -E echo "Completed GoLang Compile"
        DEPENDS golang_wolwebservice_build.out)