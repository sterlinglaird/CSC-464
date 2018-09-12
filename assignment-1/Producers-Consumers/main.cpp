#include <thread>
#include <iostream>
#include <random>
#include <chrono>
#include <mutex>
#include <condition_variable>

/*
TODO
Need to make cout thread safe and say which thread did what
*/

typedef int event;

//Does not work when buffer size is 1 because of the way I defined "empty" and "full"
const unsigned int BUFFER_SIZE = 20;
const unsigned int SEED = 1000;
const unsigned int MAX_WAIT_TIME = 500;

class Buffer {

	event buffer[BUFFER_SIZE];
	int startIdx = 0;
	int openIdx = 0;

	std::mutex mutex;
	std::condition_variable condProducer;
	std::condition_variable condConsumer;
	
public:
	event getEvent() {
		std::unique_lock<std::mutex> lock{ mutex };

		while (isEmpty()) {
			condConsumer.wait(lock);
		}

		event ev = buffer[startIdx];

		startIdx = (startIdx + 1) % BUFFER_SIZE;
		condProducer.notify_one();
		return ev;
	}

	void addEvent(event ev) {
		std::unique_lock<std::mutex> lock{ mutex };

		while (isFull()){
			condProducer.wait(lock);
		}

		buffer[openIdx] = ev;
		openIdx = (openIdx + 1) % BUFFER_SIZE;

		condConsumer.notify_one();
	}

	bool isEmpty() {
		return startIdx == openIdx;
	}

	bool isFull() {
		return startIdx == (openIdx + 1) % BUFFER_SIZE;
	}
};

Buffer buffer;

event waitForEvent() {
	int randomNum = std::rand();
	int ms = randomNum % MAX_WAIT_TIME;
	std::this_thread::sleep_for(std::chrono::milliseconds(ms));

	return event(randomNum);
}

void consumeEvent(event ev) {
	int randomNum = std::rand();
	int ms = randomNum % MAX_WAIT_TIME;
	std::this_thread::sleep_for(std::chrono::milliseconds(ms));
}

void producer() {
	while (true) {
		event ev = waitForEvent();
		std::cout << "produced: " << ev << std::endl;
		buffer.addEvent(ev);
	}
}

void consumer() {
	while (true) {
		event ev = buffer.getEvent();
		consumeEvent(ev);
		std::cout << "consumed: " << ev << std::endl;
	}
}

int main() {
	//This doesnt work for threads
	std::srand(SEED);

	std::thread producer(producer);
	std::thread consumer(consumer);

	while (true) {}

	return 0;
}