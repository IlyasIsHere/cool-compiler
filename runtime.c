// COOL Runtime Support
// This file provides the implementation of COOL runtime support functions
// that are called from the generated LLVM IR code.

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Memory management functions

// Allocate memory for an object
void* malloc_object(unsigned long size) {
    void* ptr = malloc(size);
    if (!ptr) {
        fprintf(stderr, "Runtime error: Out of memory\n");
        exit(1);
    }
    // Zero-initialize memory
    memset(ptr, 0, size);
    return ptr;
}

// IO functions

// Print a string to stdout
int out_string(const char* str) {
    return printf("%s", str);
}

// Print an integer to stdout
int out_int(int num) {
    return printf("%d", num);
}

// Read a string from stdin
char* in_string() {
    char buffer[1024];
    if (fgets(buffer, sizeof(buffer), stdin) != NULL) {
        // Remove newline if present
        size_t len = strlen(buffer);
        if (len > 0 && buffer[len-1] == '\n') {
            buffer[len-1] = '\0';
        }
        
        // Allocate and copy the string
        char* result = (char*)malloc(len + 1);
        if (result == NULL) {
            fprintf(stderr, "Runtime error: Out of memory\n");
            exit(1);
        }
        strcpy(result, buffer);
        return result;
    }
    
    // Return empty string on error
    char* result = (char*)malloc(1);
    if (result == NULL) {
        fprintf(stderr, "Runtime error: Out of memory\n");
        exit(1);
    }
    result[0] = '\0';
    return result;
}

// Read an integer from stdin
int in_int() {
    int result;
    if (scanf("%d", &result) != 1) {
        // On error, return 0
        return 0;
    }
    return result;
}

// String manipulation functions

// Get the length of a string
int string_length(const char* str) {
    return (int)strlen(str);
}

// Concatenate two strings
char* string_concat(const char* str1, const char* str2) {
    size_t len1 = strlen(str1);
    size_t len2 = strlen(str2);
    char* result = (char*)malloc(len1 + len2 + 1);
    if (result == NULL) {
        fprintf(stderr, "Runtime error: Out of memory\n");
        exit(1);
    }
    strcpy(result, str1);
    strcat(result, str2);
    return result;
}

// Get a substring
char* string_substr(const char* str, int start, int length) {
    size_t str_len = strlen(str);
    
    // Bounds checking
    if (start < 0 || start >= (int)str_len || length < 0) {
        fprintf(stderr, "Runtime error: Substring out of range\n");
        exit(1);
    }
    
    // Adjust length if necessary
    if (start + length > (int)str_len) {
        length = (int)str_len - start;
    }
    
    // Allocate memory for the substring
    char* result = (char*)malloc(length + 1);
    if (result == NULL) {
        fprintf(stderr, "Runtime error: Out of memory\n");
        exit(1);
    }
    
    // Copy the substring
    strncpy(result, str + start, length);
    result[length] = '\0';
    
    return result;
}

// Runtime support functions

// Terminate the program with an error message
void abort() {
    fprintf(stderr, "COOL program aborted\n");
    exit(1);
}

// Get the name of an object's type as a string
char* type_name(void* obj) {
    // In a real implementation, we would extract the type info from the object
    // For now, just return a placeholder
    char* result = (char*)malloc(8);
    if (result == NULL) {
        fprintf(stderr, "Runtime error: Out of memory\n");
        exit(1);
    }
    strcpy(result, "Object");
    return result;
}

// Create a shallow copy of an object
void* object_copy(void* obj) {
    // In a real implementation, we would use the object's size from its type info
    // For now, just warn that this is unimplemented
    fprintf(stderr, "Warning: object_copy not fully implemented\n");
    return obj;
}

// Case expression runtime support
void case_abort() {
    fprintf(stderr, "Runtime error: Case match failed\n");
    exit(1);
}

// Dispatch on void check
void dispatch_abort() {
    fprintf(stderr, "Runtime error: Dispatch on void\n");
    exit(1);
} 