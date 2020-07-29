# The GoImplement Encryption Library

The GoImplement encryption library is Superytpe's first client library for its proxy re-encryption implementation

## What is proxy re-encryption?

Proxy re-encryption is a new(er) cryptographic primative. It's a type of public key encryption, with the added benefit of 1-N encryption, compared to standard public key encryption's 1-1 behavior. This updated behavior is at the core of Supertype.

**Traditional public key encryption:** Alice and Bob both have two keys, one public and one secret key. If Alice wants to share an encrypted message with Bob, she can get Bob's public key and encrypt her data with it, sending it over to Bob encrypted with his public key. Bob is the only person who can decrypt this data, as his secret key is the only thing that can decrypt data encrypted with his public key.

**Proxy re-encryption:** Alice encrypts her data using _her_ public key. Separately, she creates a re-encrypiton key that allows her encrypted data to be decryptable by only her secret key, to only decryptable by Bob's secret key. This encrypted data as well as the re-encryption key is sent to an unbiased proxy (we'll call him Carter), which re-encrypts the data without every accessing the underlying plaintext. Bob then requests this now re-encrypted data from the proxy, and can now decrypt it using his secret key.

## Where does Supertype come into play?

In the above example, Alice still only needs to encrypt her data once using her own public key, as opposed to encrypting the data each time Bob, Doug, or Emily requests data with Bob's, Doug's, or Emily's public key. However, in the above example, Alice still needs to create and manage the re-encryption keys between her and Bob, Doug, or Emily - for whenever the data is requested. Not much of an improvement.

Supertype is, at its core, a managed service for proxy re-encryption. Supertype will manage an up-to-date registry of re-encryption keys between Alice and anyone else using Supertype. Similarly, Supertype will manage the encrypted data, transporting it from producer to consumer. Acting as this transport layer, all Alice needs to do is produce any new data to Supertype, and all Bob, Doug, or Emily needs to do is consume the data from Supertype.

Supertype takes care of the proxy step by effectively making everyone using the Supertype platform a proxy as well. When conusming data, using the `GoImplement` library, they will re-encrypt the data within their architecture so that they can then decrypt the data within their architecture.

Supertype provides the efficiency of symmetric encryption with the security of public key encryption, enabling the necessary amount of interconnectivity for the true Internet of Things.