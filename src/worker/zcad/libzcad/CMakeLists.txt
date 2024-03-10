cmake_minimum_required(VERSION 3.10)
project(Hello)

# Set C++ standard
set(CMAKE_CXX_STANDARD 11)

# Create a shared library
add_library(hello SHARED hello.h hello.c)
