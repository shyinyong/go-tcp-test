package handler

import (
	"errors"
	"io"
	"net"
)

// ReadMessage reads a protobuf message from the connection.
func readMessage(conn net.Conn) ([]byte, error) {
	// Create a buffer to store the incoming data
	buffer := make([]byte, 0)

	// Set a maximum buffer size to prevent large messages from causing memory issues
	maxBufferSize := 4096 // You can adjust this size based on your requirements

	// Read data from the connection in chunks until the entire message is received
	for {
		// Read a chunk of data from the connection
		chunk := make([]byte, 1024) // Adjust the chunk size as needed
		n, err := conn.Read(chunk)
		if err != nil {
			if err == io.EOF {
				// If the connection is closed by the client, return an error
				return nil, errors.New("client connection closed")
			}
			return nil, err
		}

		// Append the received chunk to the buffer
		buffer = append(buffer, chunk[:n]...)

		// Check if the message size exceeds the maximum buffer size
		if len(buffer) > maxBufferSize {
			return nil, errors.New("message size exceeds maximum buffer size")
		}

		// Check if the entire message has been received
		if n < len(chunk) {
			break
		}
	}

	return buffer, nil
}
