---
title: Convert procedural language program to Object Oriented
date: 2019-08-16
for: Myself
---

Basic building blocks of OO programs are classes. And classes are made up of - **State & Behaviour**. In procedural language, **parameters are state** & **functions are their behaviour**

<pre>
Consider below C prorgam
```c
struct bill {
    *some fields*
}
struct payment {
    *some fields
}
main() {
   struct bill food_bill;
   struct payment cash_payment;
   process_payment( food_bill, cash_payment);
   process_invoice( food_bill, cash_payment);
}
```
</pre>

*food_bill*, *cash_payment* is state & process_payment and process_invoice are behaviour.
let's encapuslate all together in a class.

<pre>
```java
class Transaction {
    bill food_bill;
    payment online_payment;
    process_payment() {
    }
    process_invoice() {       
    }     
}
   public class SomePublicClass {
       public static void main() {
           Transaction t = new Transaction(bill, payment);
           t.process_payment();
           t.process_invoice();
       }
   }
```
</pre>

Now lets see process_payment's implementation.

<pre>
```c
void process_payment (bill, payment) {
    struct paypalrequst ppeq; 
    ppreq.xx = yy;
    callToPayPal(ppeq, bill);
    logAcceptedPaymentDetails(ppreq);
}
```
</pre>

logAcceptedPaymentDetails is behviour on paypalrequst while callToPayPal is behviour of ppeq and bill. 
logAcceptedPaymentDetails() demands for PaymentGateway as separate class. Placing callToPayPal() is however tricky. But, since there is no other method with same list of parameters. So there is no point in creating separate encapsulation for bill & ppreq.

<pre>
```java
class Transaction {
    bill food_bill;
    payment cash_payment;
    process_payment () {
        PaymentGateway pg = new PaymentGateway(yy); 
        pg.callToPayPal(this);
        pg.logAcceptedPaymentDetails()        
    }
    process_invoice() {       
    }     
}
   
   class PaymentGateway {
       private int xx;
       PaymentGateway(xx){

       }
       callToPayPal(Transaction t) {
           *refers transactionId*
       }
       logAcceptedPaymentDetails() {

       }
   }
```
</pre>