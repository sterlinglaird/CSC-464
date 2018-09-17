#include <thread>
#include <iostream>
#include <random>
#include <chrono>
#include <mutex>
#include <condition_variable>

typedef int event;

//Does not work when buffer size is 1 because of the way I defined "empty" and "full"
const unsigned int BUFFER_SIZE = 20;
const unsigned int MAX_WAIT_TIME = 50;
const unsigned int NUM_PRODUCERS = 2;
const unsigned int NUM_CONSUMERS = 10;

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
std::mutex coutMutex;

event waitForEvent() {
	int randomNum = std::rand();
	int ms = randomNum % MAX_WAIT_TIME;
	std::this_thread::sleep_for(std::chrono::milliseconds(ms));

	return event(randomNum % 10000);
}

void consumeEvent(event ev) {
	int randomNum = std::rand();
	int ms = randomNum % MAX_WAIT_TIME;
	std::this_thread::sleep_for(std::chrono::milliseconds(ms));
}

void producer(int id) {
	while (true) {
		event ev = waitForEvent();

		std::unique_lock<std::mutex> lock{ coutMutex };
		std::cout << "produced: " << ev << " by: " << id << std::endl;
		lock.unlock();

		buffer.addEvent(ev);
	}
}

void consumer(int id) {
	while (true) {
		event ev = buffer.getEvent();
		consumeEvent(ev);

		std::unique_lock<std::mutex> lock{ coutMutex };
		std::cout << "consumed: " << ev << " by: " << id << std::endl;
		lock.unlock();
	}
}

int main() {
	for (int i = 0; i < NUM_PRODUCERS; i++) {
		std::thread producer(producer, i);
		producer.detach();
	}

	for (int i = 0; i < NUM_CONSUMERS; i++) {
		std::thread consumer(consumer, i);
		consumer.detach();
	}

	while (true) {}

	return 0;
}