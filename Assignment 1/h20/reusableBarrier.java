package h2o;

import java.util.concurrent.Semaphore;

public class reusableBarrier {
	int count;
	int n;
	private Semaphore mutex =new Semaphore(1);
	private Semaphore turnslite1 =new Semaphore(0);
	private Semaphore turnslite2 =new Semaphore(1);
	
	public reusableBarrier(int n) {
		this.n= n;
	}
	
	public void mywait() {
		try {
			mutex.acquire();
			count ++;
			if (count ==n) {
				turnslite2.acquire();
				turnslite1.release();
			}
			mutex.release();
			
			turnslite1.acquire();
			turnslite1.release();
			
			mutex.acquire();
			count --;
			if (count == 0) {
				turnslite1.acquire();
				turnslite2.release();
			}
			mutex.release();
			
			turnslite2.acquire();
			turnslite2.release();
		} catch (InterruptedException e) {
			e.printStackTrace();
		}
	}
}
