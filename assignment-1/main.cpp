#include <thread>
#include <iostream>
#include <random>
#include <chrono>
#include <mutex>
#include <condition_variable>

typedef int event;

const unsigned int BUFFER_SIZE = 20;
const unsigned int SEED = 1000;

class Buffer {

	event buffer[BUFFER_SIZE];
	int startIdx = 0;
	int openIdx = 0;

	std::mutex mutex;
	std::unique_lock<std::mutex> lock;
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
		lock.unlock();
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

		lock.unlock();
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
	int ms = randomNum % 500;
	std::this_thread::sleep_for(std::chrono::milliseconds(ms));

	return event(randomNum);
}

void consumeEvent(event ev) {
	int randomNum = std::rand();
	int ms = randomNum % 500;
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
	std::srand(SEED);

	std::thread producer(producer);
	std::thread consumer(consumer);

	while (true) {}

	return 0;
}