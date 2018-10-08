package savages;

import java.util.concurrent.Semaphore;

import javax.swing.plaf.basic.BasicInternalFrameTitlePane.MaximizeAction;


public class savages {
	static Semaphore mutex = new Semaphore(1);
	static Semaphore emptyPot = new Semaphore(0);
	static Semaphore fullPot = new Semaphore(0);
	static int servings = 10;
	static final int Savages_Num = 1000000;
	
	private static void putServingsInPot () {
		System.out.println("The Cook has put 10 servings in pot.");
	}
	
	private static void eat (String s) {
		System.out.println(s+" is eating");
	}
	static class Savage implements Runnable{
		@Override
		public void run()  {
			try {
				for (int i = 1; i <Savages_Num+1; i++) {
					mutex.acquire();
					if (servings == 0) {
						emptyPot.release();
						fullPot.acquire();
						servings = 10;
					}
					servings --;
					eat("Savage "+ i);
					mutex.release();
				}
			} catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}
	static class Cook implements Runnable{
		@Override
		public void run()  {
			try {
				for (int i = 0; i < Savages_Num/10 -1; i++) {
					emptyPot.acquire();
					putServingsInPot();
					fullPot.release();
				}
			} catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}
	
	public static void main (String[] args) {
		Thread[] threads = new Thread[2];
		threads[0] = new Thread(new Cook());
		threads[1] = new Thread(new Savage());
		
		long startTime = System.currentTimeMillis(); 
		for(int i = 0; i < threads.length; i++){
				  threads[i].start();
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
		/*while (true) {
			if (finish) {
				long endTime=System.currentTimeMillis(); //get end time
				System.out.println("Time Cost "+(endTime-startTime)+"ms");	
				break;
			}
		}*/
	}
}
