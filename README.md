# Openai-Forward-Fiber

Openai-Forward-Fiber is a project written in Go that provides a forwarding service for OpenAI API using Fiber. Users can utilize this service by using the generated API key, without having direct access to their own OpenAI API key.

## Features

The following features are supported by Openai-Forward-Fiber:

- **Sliding Window Algorithm-based User Rate Limiting**: Openai-Forward-Fiber incorporates a sliding window algorithm to control and limit the number of requests made by users. This ensures fair usage of API tokens among all users.

- **Token Usage Billing**: Openai-Forward-Fiber tracks the token usage for each user and provides billing based on the tokens consumed. This allows for accurate measurement of resource utilization and cost allocation.

- **Error Handling for Edge Cases**: The project handles various edge cases that may occur during the usage of the OpenAI API. It provides robust error handling mechanisms to ensure smooth operation even in challenging scenarios.

- **User Management APIs**: Openai-Forward-Fiber exposes a range of user management APIs, which can be accessed by navigating to `/swagger` to view the detailed documentation. These APIs facilitate easy management of users and their associated settings.

- **Flexible Configuration**: Most of the project settings can be customized by modifying the configuration file (`config.yaml`). The configuration file contains detailed explanations of each field, allowing users to tailor the project to their specific requirements.

## Installation

To install and run Openai-Forward-Fiber, follow these steps:

1. Clone the repository:

   ```
   git clone https://github.com/gtoxlili/openai-forward-fiber.git
   ```

2. Navigate to the project directory:

   ```
   cd openai-forward-fiber
   ```

3. Configure the project:

   Modify the `config.yaml` file to adjust the settings according to your needs. Refer to the comments within the file for guidance on each field.

4. Build the project:

   ```
   go build
   ```

5. Run the project:

   ```
   ./openai-forward-fiber
   ```

6. Access the APIs:

   Openai-Forward-Fiber will be running on `http://localhost:3000`. You can interact with the APIs using tools like cURL or import the provided Swagger documentation to explore and consume the APIs.

## Contributing

Contributions to Openai-Forward-Fiber are welcome! If you encounter any issues or have suggestions for improvements, please open an issue on the project's GitHub repository.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.

## Acknowledgements

Openai-Forward-Fiber is built upon the [Fiber](https://github.com/gofiber/fiber) framework and utilizes the OpenAI API for forwarding requests. Special thanks to the developers of these amazing tools and tec
ologies.

## Disclaimer

Please note that Openai-Forward-Fiber is a third-party project and is not officially affiliated with or endorsed by OpenAI. Use it responsibly and ensure compliance with the terms and conditions of the OpenAI API usage.
