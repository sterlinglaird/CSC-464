#include <thread>
#include <iostream>
#include <random>
#include <chrono>
#include <mutex>
#include <condition_variable>

const unsigned int MAX_WAIT_TIME = 5;
const unsigned int NUM_SAVAGES = 10;
const unsigned int MAX_SERVINGS = 10;

enum Serving {food};

class Pot {
	unsigned int numServings = 0;

	std::mutex mutex;
	std::condition_variable fullPot;
	std::condition_variable emptyPot;

public:
	void waitUntilEmpty() {
		std::unique_lock<std::mutex> lock{ mutex };
		while (numServings > 0) {
			emptyPot.wait(lock);
		}
	}

	void fillPot(unsigned int numServings){
		std::unique_lock<std::mutex> lock{ mutex };
		this->numServings = numServings;
		fullPot.notify_all();
	}

	Serving getServing(int idx) {
		std::unique_lock<std::mutex> lock{ mutex };
		while (numServings <= 0) {
			emptyPot.notify_one();
			fullPot.wait(lock);
		}

		numServings--;
		return Serving::food;
	}
};

Pot pot;
std::mutex coutMutex;

void eat(Serving serving) {
	int randomNum = std::rand();
	int ms = randomNum % MAX_WAIT_TIME;
	std::this_thread::sleep_for(std::chrono::milliseconds(ms));
}

void cook() {
	while (true) {
		pot.waitUntilEmpty();

		std::unique_lock<std::mutex> lock{ coutMutex };

		pot.fillPot(MAX_SERVINGS);

		std::cout << "cook refilled pot" << std::endl;
		lock.unlock();
	}
}

void savage(unsigned int idx) {
	while (true) {
		std::unique_lock<std::mutex> lock{ coutMutex };
		std::cout << "savage " << idx << " looking for food" << std::endl;
		lock.unlock();

		Serving serving = pot.getServing(idx);

		lock.lock();
		std::cout << "savage " << idx << " got food" << std::endl;
		lock.unlock();

		eat(serving);
	}

}

int main() {
	std::thread savages[NUM_SAVAGES];

	std::thread cook(cook);

	for (int i = 0; i < NUM_SAVAGES; i++) {
		savages[i] = std::thread(savage, i);
	}

	while (true) {}

	return 0;
}