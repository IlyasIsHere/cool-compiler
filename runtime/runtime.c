// runtime.c - Runtime support for the COOL compiler
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// String operations
char* string_concat(char* str1, char* str2) {
    if (!str1) return strdup(str2 ? str2 : "");
    if (!str2) return strdup(str1);
    
    size_t len1 = strlen(str1);
    size_t len2 = strlen(str2);
    char* result = (char*)malloc(len1 + len2 + 1);
    
    if (result) {
        strcpy(result, str1);
        strcat(result, str2);
    }
    
    return result;
}

char* string_substr(char* str, int start, int length) {
    if (!str) return strdup("");
    
    size_t str_len = strlen(str);
    
    // Bounds checking
    if (start < 0 || start >= str_len || length < 0) {
        fprintf(stderr, "Runtime error: substring out of range\n");
        exit(1);
    }
    
    // Adjust length if needed
    if (start + length > str_len) {
        length = str_len - start;
    }
    
    char* result = (char*)malloc(length + 1);
    if (result) {
        strncpy(result, str + start, length);
        result[length] = '\0';
    }
    
    return result;
}

int string_length(char* str) {
    return str ? strlen(str) : 0;
}

// IO operations
char* in_string() {
    char buffer[1024];
    if (fgets(buffer, sizeof(buffer), stdin)) {
        // Remove newline if present
        size_t len = strlen(buffer);
        if (len > 0 && buffer[len-1] == '\n') {
            buffer[len-1] = '\0';
        }
        return strdup(buffer);
    }
    return strdup("");
}

int in_int() {
    char buffer[64];
    if (fgets(buffer, sizeof(buffer), stdin)) {
        return atoi(buffer);
    }
    return 0;
}

// Other required runtime functions
void abort() {
    fprintf(stderr, "Program aborted\n");
    exit(1);
}

char* type_name(void* obj) {
    // This is a placeholder - your actual implementation will depend on
    // how you represent runtime type information
    return strdup("Object");
}

void* object_copy(void* obj) {
    // This is a placeholder - actual implementation depends on how you represent objects
    return obj;
}

void case_abort() {
    fprintf(stderr, "COOL runtime error: Case does not match any branch\n");
    exit(1);
}

void dispatch_abort() {
    fprintf(stderr, "COOL runtime error: Dispatch to void\n");
    exit(1);
} 