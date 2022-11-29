cmake_minimum_required(VERSION 3.0.0)
project(WOLWebService VERSION 0.1.0)

set(CPACK_PROJECT_NAME ${PROJECT_NAME})
set(CPACK_PROJECT_VERSION ${PROJECT_VERSION})
include(CPack)


execute_process(COMMAND uname OUTPUT_VARIABLE uname)
if (uname MATCHES "^MSYS" OR uname MATCHES "^MINGW")
    message("This is MSYS build")
    add_subdirectory("efidll")
endif()
add_subdirectory("server")