package h2o;

import java.util.Random;
import java.util.concurrent.Semaphore;

public class h20 {
	static Semaphore mutex = new Semaphore(1);
	static Semaphore oxygenQueue = new Semaphore(0);
	static Semaphore hydrogenQueue = new Semaphore(0);
	static reusableBarrier Barrier = new reusableBarrier(3);
	volatile static int oxygen = 0;
	volatile static int hydrogen = 0;
	volatile static int waterId = 1;
	public static final int water_num = 10000;

	static class Oxygen implements Runnable{
		@Override
		public void run()  {
			try {
				mutex.acquire();
				oxygen++;
				if (hydrogen >= 2) {
					hydrogenQueue.release(2);
					hydrogen -=2;
					oxygenQueue.release();
					oxygen--;
				} else {
					mutex.release();
				}
				oxygenQueue.acquire();
				bonder.bond(0);
				Barrier.mywait();
				mutex.release();
			} catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}
	
	static class Hydrogen implements Runnable {
		@Override
		public void run() {
			try {
				mutex.acquire();
				hydrogen++;
				if (hydrogen >= 2 && oxygen >= 1) {
					hydrogenQueue.release(2);
					hydrogen -=2;
					oxygenQueue.release();
					oxygen--;
				} else {
					mutex.release();
				}
				hydrogenQueue.acquire();
				bonder.bond(1);
				Barrier.mywait();
			}  catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}
	
	static class bonder {//0:Oxygen, 1:Hydrogen
		static int bonded = 0;
		static int oxygen_num = 0;
		static int hydropgen_num = 0;
		public static synchronized void bond(int flag) {
			if (flag == 0) {
				oxygen_num++;
			} else if (flag == 1){
				hydropgen_num++;
			} else {
				System.out.println("bond error - invalid flag");
			}
			if (oxygen_num == 1 && hydropgen_num == 2) {
				oxygen_num = 0;
				hydropgen_num = 0;
				bonded++;
				System.out.printf("%d Water Bonded\n",bonded);
			} else if (oxygen_num > 1 || hydropgen_num > 2) {
				System.out.println("bond error - too many component");
			}
		}
	}
	public static void main (String[] args) {
		Oxygen oxygen = new Oxygen();
		Hydrogen hydrogen = new Hydrogen();
		Thread[] threads = new Thread[water_num*3];
		Random rand = new Random();
		int oxygen_num = 0;
		int hydrogen_num = 0;
		ThreadGroup tg = new ThreadGroup("total");
		long startTime = System.currentTimeMillis(); 
		for (int i=0; i<water_num*3;i++) {
			int ran = rand.nextInt(3);
			if ((ran == 0 && oxygen_num <water_num) || hydrogen_num == water_num*2) {
				Thread oxygenThread = new Thread(oxygen);
				threads[i] = oxygenThread;
				oxygenThread.start();
				oxygen_num++;
			} else {
				Thread hydrogenThread = new Thread(hydrogen);
				threads[i] = hydrogenThread;
				hydrogenThread.start();
				hydrogen_num++;
			}
		}
		
		for(int j = 0; j < threads.length; j++){
			  try {
				  threads[j].join();
			  } catch (Exception e) {
				  e.printStackTrace();
			  }
		}
		long endTime=System.currentTimeMillis(); //get end time
		System.out.println("Time Cost "+(endTime-startTime)+"ms");
	}
}
