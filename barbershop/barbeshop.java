package barbeshop;

import java.util.concurrent.Semaphore;


public class barbeshop {
	static Semaphore mutex = new Semaphore(1);
	static Semaphore customer = new Semaphore(0);
	static Semaphore barber = new Semaphore(0);
	static Semaphore customerDone = new Semaphore(0);
	static Semaphore barberDone = new Semaphore(0);
	static final int n = 4;
	static int customers = 0;
	static final int Customers_Num = 10000;
	
	
	
	private static void balk(String s) {
		System.out.println(s+" left because Barbershop is full");
	}
	
	private static void getHairCut(String s) {
		try {
			Thread.sleep(15);
		} catch (InterruptedException e){
			e.printStackTrace();
		}
		System.out.println(s+" is getting hair cut ");
	}
	
	private static void cutHair() {
		try {
			Thread.sleep(100000);
		} catch (InterruptedException e){
			e.printStackTrace();
		}
		//System.out.println("Barber is working.");
	}
	
	static class Customer implements Runnable{
		int id;
		
		public void setId(int id) {
			this.id = id;
		}
		
		@Override
		public void run()  {
			try {
				mutex.acquire();
				if (customers == n) {
					balk("Customer "+id);
					mutex.release();
				} else {
					customers++;
					mutex.release();
						
					customer.release();
					barber.acquire();
						
					getHairCut("Customer "+id);
						
					customerDone.release();
					barberDone.acquire();
						
					mutex.acquire();
					customers--;
					mutex.release();
				}
			} catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}
	
	static class Barber implements Runnable{
		@Override
		public void run()  {
			try {
				while (true) {
					customer.acquire();
					barber.release();
					
					//cutHair();
					
					customerDone.acquire();
					barberDone.release();
				}
			} catch (InterruptedException e){
				e.printStackTrace();
			}
	    }
	}
	
	public static void main (String[] args) {
		Thread[] threads = new Thread[Customers_Num + 1];
		threads[0] = new Thread(new Barber());
		for(int i = 1; i < threads.length; i++){
			Customer c = new Customer();
			c.setId(i);
			threads[i] = new Thread(c);
		}
		long startTime = System.currentTimeMillis(); 
		long beforeUsedMem=Runtime.getRuntime().totalMemory()-Runtime.getRuntime().freeMemory();
		for(int i = 0; i < threads.length; i++){
				  threads[i].start();
				  try {
					Thread.sleep(10);
				  } catch (InterruptedException e){
					e.printStackTrace();
				  }
		}
		for(int j = 1; j < threads.length; j++){
			  try {
				  threads[j].join();
			  } catch (Exception e) {
				  e.printStackTrace();
			  }
		}
		long afterUsedMem=Runtime.getRuntime().totalMemory()-Runtime.getRuntime().freeMemory();
		long endTime=System.currentTimeMillis(); //get end time
		System.out.println("Time Cost "+(endTime-startTime)+"ms");	
		System.out.println("Memory Usage :"+(afterUsedMem-beforeUsedMem)/(1024.0 * 1024.0)+" MB");	

		
		
	}
}
