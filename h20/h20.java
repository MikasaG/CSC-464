package h2o;

import java.util.Random;
import java.util.concurrent.Semaphore;

public class h20 {
	static Semaphore mutex = new Semaphore(1);
	static Semaphore oxygenQueue = new Semaphore(0);
	static Semaphore hydrogenQueue = new Semaphore(0);
	static Semaphore barrier = new Semaphore(3);
	static int oxygen = 0;
	static int hydrogen = 0;
	static int waterId = 1;
	public static final int water_num = 9;

	static class Oxygen extends Thread {
		@Override
		public void run()  {
			try {
				mutex.acquire();
			} catch (InterruptedException e){
				e.printStackTrace();
			}
			oxygen++;
			if (hydrogen >= 2) {
				hydrogenQueue.release(2);
				hydrogen -=2;
				oxygenQueue.release();
				oxygen--;
			} else {
				mutex.release();
			}
			try {
				oxygenQueue.acquire();
			} catch (InterruptedException e){
				e.printStackTrace();
			}
			bond();
			try {
				barrier.acquire();
			} catch (InterruptedException e){
				e.printStackTrace();
			}
			mutex.release();
	    }
	}
	
	static class Hydrogen extends Thread {
		@Override
		public void run() {
			try {
				mutex.acquire();
			} catch (InterruptedException e){
				e.printStackTrace();
			}
			hydrogen++;
			if (hydrogen >= 2 && oxygen >= 1) {
				hydrogenQueue.release(2);
				hydrogen -=2;
				oxygenQueue.release();
				oxygen--;
			} else {
				mutex.release();
			}
			try {
				hydrogenQueue.acquire();
			} catch (InterruptedException e){
				e.printStackTrace();
			}
			bond();
			try {
				barrier.acquire();
			} catch (InterruptedException e){
				e.printStackTrace();
			}
			mutex.release();
	    }
	}
	
	public static synchronized void bond() {
		System.out.printf("%d Water Bonded\n",waterId);
		waterId++;
	}
	public static void main (String[] args) {
		Random rand = new Random();
		int oxygen_num = 0;
		int hydrogen_num = 0;
		for (int i=0; i<water_num*3;i++) {
			int ran = rand.nextInt(3);
			if ((ran == 0 && oxygen_num <water_num) || hydrogen_num == water_num*2) {
				Oxygen oxygenThread = new Oxygen();
				oxygenThread.start();
				oxygen_num++;
				//System.out.printf("%d Oxygen generate\n",oxygen_num);
			} else {
				Hydrogen hydrogenThread = new Hydrogen();
				hydrogenThread.start();
				hydrogen_num++;
				//System.out.printf("%d Hydrogen generate\n",hydrogen_num);
			}
		}
	}
}
