cmake_minimum_required(VERSION 3.0.0)
add_definitions(-DUNICODE -D_UNICODE)
FILE(GLOB MyHeaders ${CMAKE_CURRENT_SOURCE_DIR}/*.h)
FILE(GLOB MySources ${CMAKE_CURRENT_SOURCE_DIR}/*.cpp)
# Create a target for the library
ADD_LIBRARY(efiDLL SHARED ${MySources} ${MyHeaders})
target_compile_definitions(efiDLL PRIVATE _WIN32_WINNT=0x0A00 WINVER=0x0A00)
target_link_options(efiDLL PRIVATE -static)
set_target_properties(efiDLL 
                    PROPERTIES PREFIX ""
                    LINKER_LANGUAGE CXX
                    LIBRARY_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/dist"
                    RUNTIME_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/dist")


