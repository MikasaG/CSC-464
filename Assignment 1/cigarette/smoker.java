package smoker;

import java.util.concurrent.Semaphore;


public class smoker {
	static Semaphore mutex = new Semaphore(1);
	static Semaphore tobacco = new Semaphore(0);
	static Semaphore paper = new Semaphore(0);
	static Semaphore match = new Semaphore(0);
	static Semaphore tobaccoSem = new Semaphore(0);
	static Semaphore paperSem = new Semaphore(0);
	static Semaphore matchSem = new Semaphore(0);
	static Semaphore threads_sem = new Semaphore(0);
	static Semaphore agent_sem = new Semaphore(1);
	static boolean isTobacco = false;
	static boolean isPaper = false;
	static boolean isMatch = false;
	static final int MAX_Cigarette =10000; 
	static int Curresnt_Cigatette = 0;
	volatile static boolean finish = false;
	
	private static void makeCigarette(String s) {
		System.out.println(s + " has made a cigarette");
	}
	
	private static void smoke(String s) {
		System.out.println(s + " is smoking");
	}
	
	static class Smoker_WithTobacco implements Runnable{
		@Override
		public void run()  {
			try {
				while (!finish) {
					tobaccoSem.acquire();
					makeCigarette("Smoker A");
					agent_sem.release();
					smoke("Smoker A");
				}
			} catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}
	static class Smoker_WithPaper implements Runnable{
		@Override
		public void run()  {
			try {
				while (!finish) {
					paperSem.acquire();
					makeCigarette("Smoker B");
					agent_sem.release();
					smoke("Smoker B");
				}
			} catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}
	static class Smoker_WithMatch implements Runnable{
		@Override
		public void run()  {
			try {
				while (!finish) {
					matchSem.acquire();
					makeCigarette("Smoker C");
					agent_sem.release();
					smoke("Smoker C");
				}
			} catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}
	static class Pusher_A implements Runnable{//check for tobacco
		@Override
		public void run()  {
			try {
				while (!finish) {
					tobacco.acquire();
					mutex.acquire();
					if (isPaper) {
						isPaper = false;
						matchSem.release();
					} else if (isMatch) {    
						isMatch = false;
						paperSem.release();
					} else {
						isTobacco = true;
					}
					mutex.release();
				}
			} catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}
	static class Pusher_B implements Runnable{//check for paper
		@Override
		public void run()  {
			try {
				while (!finish) {
					paper.acquire();
					mutex.acquire();
					if (isTobacco) {
						isTobacco = false;
						matchSem.release();
					} else if (isMatch) {    
						isMatch = false;
						tobaccoSem.release();
					} else {
						isPaper = true;
					}
					mutex.release();
				}
			} catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}
	static class Pusher_C implements Runnable{//check for match
		@Override
		public void run()  {
			try {
				while (!finish) {
					match.acquire();
					mutex.acquire();
					if (isPaper) {
						isPaper = false;
						tobaccoSem.release();
					} else if (isTobacco) {    
						isTobacco = false;
						paperSem.release();
					} else {
						isMatch = true;
					}
					mutex.release();
				}
			} catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}
	static class Agent_A implements Runnable{
		@Override
		public void run()  {
			try {
				while (!finish) {
					agent_sem.acquire();
					mutex.acquire();
					tobacco.release();
					paper.release();
					Curresnt_Cigatette++;
					if (Curresnt_Cigatette > MAX_Cigarette) {
						finish = true;
					}
					mutex.release();
					System.out.println("Agent has put tobacco and paper on Table");
				}
			} catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}
	static class Agent_B implements Runnable{
		@Override
		public void run()  {
			try {
				while (!finish) {
					agent_sem.acquire();
					mutex.acquire();
					paper.release();
					match.release();
					Curresnt_Cigatette++;
					if (Curresnt_Cigatette > MAX_Cigarette) {
						finish = true;
					}
					mutex.release();
					System.out.println("Agent has put paper and match on Table");
				}
			} catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}
	static class Agent_C implements Runnable{
		@Override
		public void run()  {
			try {
				while (!finish) {
					agent_sem.acquire();
					mutex.acquire();
					match.release();
					tobacco.release();
					Curresnt_Cigatette++;
					if (Curresnt_Cigatette > MAX_Cigarette) {
						finish = true;
					}
					mutex.release();
					System.out.println("Agent has put match and tobacco on Table");
				}
			} catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}
	
	public static void main (String[] args) {
		Thread[] threads = new Thread[9];
		threads[0] = new Thread(new Agent_A());
		threads[1] = new Thread(new Agent_B());
		threads[2] = new Thread(new Agent_C());
		threads[3] = new Thread(new Pusher_A());
		threads[4] = new Thread(new Pusher_B());
		threads[5] = new Thread(new Pusher_C());
		threads[6] = new Thread(new Smoker_WithTobacco());
		threads[7] = new Thread(new Smoker_WithPaper());
		threads[8] = new Thread(new Smoker_WithMatch());
		long startTime = System.currentTimeMillis(); 
		long beforeUsedMem=Runtime.getRuntime().totalMemory()-Runtime.getRuntime().freeMemory();
		for(int i = 0; i < threads.length; i++){
				  threads[i].start();
		}
		/*for(int j = 0; j < threads.length; j++){
			  try {
				  threads[j].join();
			  } catch (Exception e) {
				  e.printStackTrace();
			  }
		}*/
		while (true) {
			if (finish) {
				long afterUsedMem=Runtime.getRuntime().totalMemory()-Runtime.getRuntime().freeMemory();
				long endTime=System.currentTimeMillis(); //get end time
				System.out.println("Time Cost "+(endTime-startTime)+"ms");	
				System.out.println("Memory Usage :"+(afterUsedMem-beforeUsedMem)/(1024.0 * 1024.0)+" MB");	
				break;
			}
		}
		
	}
}
